package main

func (exp *Expression) GetVars() map[Var]bool {
	return exp.vars
}
