package forms

type errors map[string][]string

//ADD MESSAGE FOR A GIVEN FORM FIELD
func (e errors) Add(field, message string){
	e[field] = append(e[field], message)
}
//RETURNS THE 1ST ERROR MESSAGE
func (e errors) Get(field string) string{
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}