package util

// CalcularDigitoModulo11 genera el dígito verificador para la clave de acceso.
func CalcularDigitoModulo11(clave string) int {
	suma := 0
	factor := 2
	for i := len(clave) - 1; i >= 0; i-- {
		digito := int(clave[i] - '0')
		suma += digito * factor
		factor++
		if factor > 7 {
			factor = 2
		}
	}

	residuo := suma % 11
	digitoVerificador := 11 - residuo

	if digitoVerificador == 11 {
		digitoVerificador = 0
	} else if digitoVerificador == 10 {
		digitoVerificador = 1
	}

	return digitoVerificador
}
