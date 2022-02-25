## Paper Summary

#### *Devign: Effective Vulnerability Identification by Learning Comprehensive Program Semantics via Graph Neural Networks*



#### 1. goal of the paper

Vulnerabilities identification is a crucial area protecting software system from being attacked. However, the existing vulnerability analysis methods are limited to cover various vulnerabilities. And it is inefficient to identify vulnerabilities manually. 

This paper tries to implement a more efficient and accurate methodology to automatedly analyze and identify the potential vulnerabilities from the source code. 



#### 2. prior solutions

static analysis,  dynamic analysis, symbolic execution

using machine learning 

​		-- treat source code as flat natural language sequences as input. 

​		-- probe more structural neural networks, like tree and graph structures



#### 3. main technology

Devign algorithm:

1. Graph embedding layer -- to overcome the limitations of expressing logic and structures in the source code, this paper proposes the composite code representation, combining and encoding AST, CFG, DFG and NCS into a joint graph to remain the semantics of source code. 

2. Gated graph recurrent layer -- learning features of nodes by aggregating and transferring information from neighboring nodes in the graph
3.  Conv module -- selecting meaningful node representations for graph level classification

#### 4. evaluation and results

1. collect labeled data sets from 4 popular C libraries manually. And evaluate the effectiveness through these data sets. 

   

2. compare the results with other learning based module like CNN, 3-layer BiLSTM, etc.

   ​	-- based on the ACC and F1-score, Devign has 10% improvment

3. compare the results between Devign with Conv module and Ggrn with flat summation

   ​	-- Conv powered Devign has 6% improvement

4. compare the results of Devign on different types of code representations (single and complex code representations , subgraph and composite graph)

   ​	-- For the Ggrn model, there is little difference between the composite graph and the single graph, while for the Devign model, the composite graph is better than the single graph

5. compare results between Devign method and classic methods like Cppcheck on non-balanced data sets

   ​	-- Devign's F1-scores are 20-30% greater than the classic methods. 

6. test and verify if Devign could work on recent vulnerabilities

   ​	-- Devign has 74% accuracy on identifying 0-day vulnerabilities collected from 40 CVEs.  

#### 5. limitations

1. one of the core technology in Devign is the composite graph. However, the results using composite graphs are just a little better than the results using single graphs.  More experiments are needed to rule out errors. 

2. graphs containing over 500 nodes are excluded, about 15% of the whole data sets. These data may influence the results. 

3. only 0-day vulnerabilities were tested by Devign. Should do verification on more vulnerabilities. 

   

#### 6. potential improvements

1. the test based on binary classification. Try to implement multiclass  classification.

2. the labeling and training are based on the function level. What about testing on other levels?

3. are there better ways to simplify and embed DFG?

   

#### 7. reproducing

authors has a github repository, but this repository only implements a AST-based version and cannot using composite graph. 









