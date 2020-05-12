# goptim

WRS (Weighted Random Search) is an improved version of Random Search (RS), used here for hyperparameter optimization of machine learning algorithms. Unlike the standard RS, which generates for each trial new values for all hyperparameters, we generate new values for each hyperparameter with a probability of change. The intuition behind our approach is that a value that already triggered a good result is a good candidate for the next step, and should be tested in new combinations of hyperparameter values. Within the same computational budget, our method yields better results than the standard RS. Our theoretical results prove this statement.

The paper describing WRS can be found at http://univagora.ro/jour/index.php/ijccc/article/view/3514

Please cite the above paper as:

```
@article{IJCCC3514,
	author = {Adrian-Catalin Florea and Razvan Andonie},
	title = {Weighted Random Search for Hyperparameter Optimization},
	journal = {International Journal of Computers Communications & Control},
	volume = {14},
	number = {2},
	year = {2019},
	keywords = {hyperparameter optimization, random search, deep learning, neural networks},
	abstract = {We introduce an improved version of Random Search (RS), used here for hyperparameter optimization of machine learning algorithms. Unlike the standard RS, which generates for each trial new values for all hyperparameters, we generate new values for each hyperparameter with a probability of change. The intuition behind our approach is that a value that already triggered a good result is a good candidate for the next step, and should be tested in new combinations of hyperparameter values. Within the same computational budget, our method yields better results than the standard RS. Our theoretical results prove this statement. We test our method on a variation of one of the most commonly used objective function for this class of problems (the Grievank function) and for the hyperparameter optimization of a deep learning CNN architecture. Our results can be generalized to any optimization problem dened on a discrete domain.},
	issn = {1841-9844},	pages = {154--169},	doi = {10.15837/ijccc.2019.2.3514},
	url = {http://univagora.ro/jour/index.php/ijccc/article/view/3514}
}
```



