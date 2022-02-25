## simple C/C++ program analyzer

the code below is based on the clang tutorial :

https://clang.llvm.org/docs/RAVFrontendAction.html

```c++
namespace clang
{
// For each source file provided to the tool, a new FrontendAction is created.
class MyFrontendAction : public ASTFrontendAction
{
public:
    MyFrontendAction() = default;
    void EndSourceFileAction() override;

    std::unique_ptr<ASTConsumer> CreateASTConsumer(CompilerInstance& CI, StringRef file) override;

private:
    Rewriter TheRewriter;
};
} // namespace clang

int Function(int argc, const char** argv);
//-------------------------------------------------------------------------
//example.cpp
//-------------------------------------------------------------------------
//------------------------------------------------------------------------------
// Tooling sample. Demonstrates:
//
// * How to write a simple source tool using libTooling.
// * How to use RecursiveASTVisitor to find interesting AST nodes.
// * How to use the Rewriter API to rewrite the source code.
//
// Eli Bendersky (eliben@gmail.com)
// This code is in the public domain
//------------------------------------------------------------------------------
#include <sstream>
#include <string>
#include <iostream>

#include "clang/AST/AST.h"
#include "clang/AST/ASTConsumer.h"
#include "clang/AST/RecursiveASTVisitor.h"
#include "clang/Frontend/ASTConsumers.h"
#include "clang/Frontend/CompilerInstance.h"
#include "clang/Tooling/CommonOptionsParser.h"
#include "clang/Tooling/Tooling.h"
#include "llvm/ADT/STLExtras.h"
#include "llvm/Support/raw_ostream.h"

int gotoCounter = 0;


using namespace clang;
using namespace clang::driver;
using namespace clang::tooling;

static llvm::cl::OptionCategory ToolingSampleCategory("Tooling Sample");

// By implementing RecursiveASTVisitor, we can specify which AST nodes
// we're interested in by overriding relevant methods.
class MyASTVisitor : public RecursiveASTVisitor<MyASTVisitor>
{
public:
    MyASTVisitor(Rewriter& R) : TheRewriter(R) {}
    
    

    bool VisitStmt(Stmt* s)
    {
    	
        if (isa<GotoStmt>(s)) 
        {
        	gotoCounter++;
    		
        	auto* GotoStatement = cast<GotoStmt>(s);
        	//std::cout  << "The stmtCount of goto stmt is :" << gotoCounter << std::endl;
        	TheRewriter.InsertText(GotoStatement->getBeginLoc(), "// the 'goto' part\n", true,true);
        	std::stringstream SStemp;
        	SStemp << "// the count of goto statement at this point is :  " << gotoCounter << "\n";
        	TheRewriter.InsertText(ForStatement->getEndLoc().getLocWithOffset(1), SStemp.str(), true,true);
        
        	return true;
        }
        
    }

    

private:
    Rewriter& TheRewriter;
};

// Implementation of the ASTConsumer interface for reading an AST produced
// by the Clang parser.
class MyASTConsumer : public ASTConsumer
{
public:
    MyASTConsumer(Rewriter& R) : Visitor(R) {}

    // Override the method that gets called for each parsed top-level
    // declaration.
    bool HandleTopLevelDecl(DeclGroupRef DR) override
    {
        for (auto& b : DR)
        {
            // Traverse the declaration using our AST visitor.
            Visitor.TraverseDecl(b);
            //b->dump();
        }
        return true;
    }

private:
    MyASTVisitor Visitor;
};

std::unique_ptr<ASTConsumer> MyFrontendAction::CreateASTConsumer(CompilerInstance& CI, StringRef file)
{
    llvm::errs() << "** Creating AST consumer for: " << file << "\n";
    TheRewriter.setSourceMgr(CI.getSourceManager(), CI.getLangOpts());
    return std::make_unique<MyASTConsumer>(TheRewriter);
}
void MyFrontendAction::EndSourceFileAction()
{
    SourceManager& SM = TheRewriter.getSourceMgr();
    llvm::errs() << "** EndSourceFileAction for: " << SM.getFileEntryForID(SM.getMainFileID())->getName() << "\n";

    // Now emit the rewritten buffer.
    //TheRewriter.getEditBuffer(SM.getMainFileID()).write(llvm::outs());
    std::error_code error_code;
    llvm::raw_fd_ostream outFile("output.cpp", error_code, llvm::sys::fs::OF_None);
    //this will write the result to outFile
    TheRewriter.getEditBuffer(SM.getMainFileID()).write(outFile);
    outFile.close();
    //TheRewriter.getEditBuffer(SM.getMainFileID()).write(llvm::outs());
    printf("there are %d Goto statements in this file \n", gotoCounter);
}

int main(int argc, const char** argv)
{
    

    // ClangTool::run accepts a FrontendActionFactory, which is then used to
    // create new objects implementing the FrontendAction interface. Here we use
    // the helper newFrontendActionFactory to create a default factory that will
    // return a new MyFrontendAction object every time.
    // To further customize this, we could create our own factory class.
    auto ExpectedParser = CommonOptionsParser::create(argc, argv, ToolingSampleCategory);
    if (!ExpectedParser) {
        // Fail gracefully for unsupported options.
        llvm::errs() << ExpectedParser.takeError();
        return 1;
    }
    CommonOptionsParser& OptionsParser = ExpectedParser.get();
    ClangTool Tool(OptionsParser.getCompilations(),
                 OptionsParser.getSourcePathList());
    return Tool.run(newFrontendActionFactory<MyFrontendAction>().get());
}

```



### Clang AST

AST (Abstract Syntax Tree) is an intermediate product in the execution process of the front-end compilation. The code is translated by the front-end of the compiler into an abstract tree-like structure that can be described in a limited language.

Clang AST does the similar things, but it has a lot of useful things compared to the AST generated by other compilers. These things can make Clang AST behave like ordinary C++ code variables, and allow manipulating AST and getting information from AST easily.



### clang tools

The sample code this time is actually a stand-alone command-line program that uses the libclangTooling library. Therefore, Clang Tools is an independent running tool implemented through the function library provided by Clang.



### RecursiveASTVisitor

When we get the AST provided by Clang, its top-level structure must be a "**TranslationUnit**". We can traverse to find a node, but we can also use the algorithm packaged by Clang. The package of this algorithm is **RecursiveASTVisitor**. We inherit this class and implement **VisitCXXRecordDecl**, then this method will be triggered on the node that visits the **CXXRecordDecl** type. (**CXXRecordDecl** type is used to represent C++ class/union/struct)

```c++
class FindNamedClassVisitor : public RecursiveASTVisitor&lt;FindNamedClassVisitor&gt; {
public:
  bool VisitCXXRecordDecl(CXXRecordDecl *Declaration) {
    // what you want to do on the selected node
    
    // This return value indicates whether you need to visit other nodes next
    return true;
  }
};
```

**RecursiveASTVisitor** traverses each **CXX** node in a depth-first search order, and calls the **VisitCXX()** method in the class template for each **CXX** node in the AST tree. If **VisitCXX** returns false, the recursive traversal will end.



### FrontendAction & ASTConsumer

We have implemented **RecursiveASTVisitor** earlier, and it is reasonable to realize the core logic of searching symbols based on Clang AST. But obviously we need an entry that defines how we get the entire AST through a compiler instance and how to access the AST.

When writing Clang-based tools (such as Clang plugins) or stand-alone tools based on LibTooling, a common entry point is **FrontendAction**. It allows user-specific actions to be performed during compilation. If you want to run a tool on a Clang AST tree, a convenient interface **ASTFrontendAction** is provided, which is responsible for performing the action. The only part left is to implement the CreateASTConsumer method, which returns an ASTConsumer for each translation unit.

```c++
class MyFrontendAction : public ASTFrontendAction
{
public:
    MyFrontendAction() = default;
    void EndSourceFileAction() override;

    std::unique_ptr<ASTConsumer> CreateASTConsumer(CompilerInstance& CI, StringRef file) override;

private:
    Rewriter TheRewriter;
};
```

The parsed source code file will be passed back through **ASTConsumer**. So we need to implement an **ASTConsumer** to obtain and traverse the AST.

```c++
// Implementation of the ASTConsumer interface for reading an AST produced by the Clang parser.
class MyASTConsumer : public ASTConsumer
{
public:
    MyASTConsumer(Rewriter& R) : Visitor(R) {}

    // Override the method that gets called for each parsed top-level declaration.
    bool HandleTopLevelDecl(DeclGroupRef DR) override
    {
        for (auto& b : DR)
        {
            // Traverse the declaration using our AST visitor.
            Visitor.TraverseDecl(b);
            //b->dump();
        }
        return true;
    }

private:
    MyASTVisitor Visitor;
};

std::unique_ptr<ASTConsumer> MyFrontendAction::CreateASTConsumer(CompilerInstance& CI, StringRef file)
{
    llvm::errs() << "** Creating AST consumer for: " << file << "\n";
    TheRewriter.setSourceMgr(CI.getSourceManager(), CI.getLangOpts());
    return std::make_unique<MyASTConsumer>(TheRewriter);
}
```

When the AST tree is built, **MyFrontendAction::CreateASTConsumer** will be called to use our custom implemented **ASTConsumer** and return the relevant nodes to us. 

A **TheRewriter** member can be found in **MyFrontendAction**, which is a rewriter, mainly used to write back the code that we added comments to the node we visit.

**HandleTopLevelDecl** will call back to us the corresponding node information, just use **MyASTVisitor** to achieve the function we want.

### EndSourceFileAction

**EndSourceFileAction** callback at the end of processing a single input.  We have added some comments using **TheRewriter**, and we want to remain it to the source code after parsing. So we emit the rewritten buffer to a new file in this method.

```c++
void MyFrontendAction::EndSourceFileAction()
{
    SourceManager& SM = TheRewriter.getSourceMgr();
    llvm::errs() << "** EndSourceFileAction for: " << SM.getFileEntryForID(SM.getMainFileID())->getName() << "\n";

    // Now emit the rewritten buffer.
    //TheRewriter.getEditBuffer(SM.getMainFileID()).write(llvm::outs());
    std::error_code error_code;
    llvm::raw_fd_ostream outFile("output.cpp", error_code, llvm::sys::fs::OF_None);
    //this will write the result to outFile
    TheRewriter.getEditBuffer(SM.getMainFileID()).write(outFile);
    outFile.close();
    //TheRewriter.getEditBuffer(SM.getMainFileID()).write(llvm::outs());
    printf("there are %d Goto statements in this file \n", gotoCounter);
}
```

 

### main function

use the interface to pass in the source file as argument for parsing.





