package payment

type WhereCondition string

func (wc WhereCondition) AddCondition(condition string) WhereCondition {
	if len(wc) == 0 {
		wc += ` WHERE `
	} else {
		wc += ` AND `
	}
	wc += WhereCondition(condition)

	return wc
}

func (wc WhereCondition) String() string {
	return string(wc)
}
