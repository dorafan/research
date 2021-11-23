i use a slide to record my work progress. 

The goal is to find an efficient methodology to find potential semantic errors in Webassembly runtimes implemented in different ways.

Webassmbly, or wasm,  only has specifications.  Developers build the runtime in their own ways. There might be many potential bugs in these implementations and even some hidden semantic, which is hard to be found using traditional testing techs like fuzzing. 

A easy way to find sementic error is to use differential testing. You should get the same results if you feed the same input to different implementations. The question is efficiency. And that is the goal of this project. 

Only do bruce fuzzing is inefficient. Coverage-guided fuzzing is a good improvement. Try to make the test cover the entire code, and focus on seldom used code branches or code that already has bugs. This oriented test method is far more efficient than bruce fuzzing. 

The former method focuses on only one target and crash error. If the target could run and exit successfully even with a wrong result, testing will pass.  However, in this project, we want to use one method to find semantic bugs in runtimes built in totally different ways.   How to combine existing methods and make progress is what we need to do. 