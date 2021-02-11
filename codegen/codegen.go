package codegen

import (
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/rumpl/monkey-lang/ast"
	"github.com/rumpl/monkey-lang/object"
	"tinygo.org/x/go-llvm"
)

type CG struct {
	program       ast.Node
	targetMachine llvm.TargetMachine
	builder       llvm.Builder
	mod           llvm.Module
}

func New(program ast.Node) *CG {
	return &CG{
		program: program,
	}
}

func (c *CG) Codegen(env *object.Environment) error {
	if err := llvm.InitializeNativeTarget(); err != nil {
		return err
	}
	if err := llvm.InitializeNativeAsmPrinter(); err != nil {
		return err
	}
	llvm.InitializeAllAsmParsers()

	target, _ := llvm.GetTargetFromTriple(llvm.DefaultTargetTriple())

	c.targetMachine = target.CreateTargetMachine(
		llvm.DefaultTargetTriple(),
		"",
		"",
		llvm.CodeGenLevelNone,
		llvm.RelocDefault,
		llvm.CodeModelDefault,
	)
	passManager := llvm.NewPassManager()
	defer passManager.Dispose()

	passManager.AddCFGSimplificationPass()
	passManager.AddConstantMergePass()
	passManager.AddGVNPass()
	passManager.AddReassociatePass()

	c.builder = llvm.NewBuilder()
	c.mod = llvm.NewModule("main")

	c.codegen(c.program, env)
	if ok := llvm.VerifyModule(c.mod, llvm.PrintMessageAction); ok != nil {
		fmt.Println(ok.Error())
	}

	passManager.Run(c.mod)
	c.mod.Dump()

	llvmBuf, _ := c.targetMachine.EmitToMemoryBuffer(c.mod, llvm.ObjectFile)
	_ = ioutil.WriteFile("out.o", llvmBuf.Bytes(), 0666)

	return exec.Command("cc", "out.o", "-fno-PIE", "-lc", "-o", "out").Run()
}

func (c *CG) codegen(node ast.Node, env *object.Environment) llvm.Value {
	switch node := node.(type) {
	case *ast.Program:
		return c.codegenProgram(node, env)
	case *ast.BlockStatement:
		return c.codegenBlockStatement(node, env)
	case *ast.FunctionLiteral:
		main := llvm.FunctionType(llvm.Int32Type(), []llvm.Type{}, false)
		v := llvm.AddFunction(c.mod, "main", main)
		block := llvm.AddBasicBlock(c.mod.NamedFunction("main"), "entry")
		c.builder.SetInsertPoint(block, block.FirstInstruction())

		c.codegen(node.Body, env)

		return v
	case *ast.ExpressionStatement:
		return c.codegen(node.Expression, env)
	case *ast.InfixExpression:
		left := c.codegen(node.Left, env)
		right := c.codegen(node.Right, env)
		return c.codegenInfixExpression(node.Operator, left, right)
	case *ast.ReturnStatement:
		val := c.codegen(node.ReturnValue, env)
		c.builder.CreateRet(val)
		return val
	case *ast.IntegerLiteral:
		b := c.builder.CreateAlloca(llvm.Int32Type(), "")
		c.builder.CreateStore(llvm.ConstInt(llvm.Int32Type(), uint64(node.Value), false), b)
		return b
	}

	return llvm.Value{}
}

func (c *CG) codegenProgram(program *ast.Program, env *object.Environment) llvm.Value {
	var result llvm.Value

	for _, stmt := range program.Statements {
		result = c.codegen(stmt, env)
	}

	return result
}

func (c *CG) codegenBlockStatement(block *ast.BlockStatement, env *object.Environment) llvm.Value {
	var result llvm.Value

	for _, stmt := range block.Statements {
		result = c.codegen(stmt, env)
	}

	return result
}

func (c *CG) codegenInfixExpression(operator string, left llvm.Value, right llvm.Value) llvm.Value {
	aVal := c.builder.CreateLoad(left, "")
	bVal := c.builder.CreateLoad(right, "")

	var result llvm.Value
	switch operator {
	case "+":
		result = c.builder.CreateAdd(aVal, bVal, "")
	case "-":
		result = c.builder.CreateSub(aVal, bVal, "")
	case "*":
		result = c.builder.CreateMul(aVal, bVal, "")
	case "/":
		result = c.builder.CreateFDiv(aVal, bVal, "")
	}

	return result
}
