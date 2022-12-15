package main

func StrategieAcheteurSimple(offre int, prixMax int, round int) (contreOffre int) {
	difference := offre - prixMax
	contreOffre = prixMax - difference
	return contreOffre
}

func StrategieVendeurSimple(offre int, offrePrecedente int, prixMin int, round int) (contreOffre int, possible bool) {
	milieu := (offrePrecedente + offre) / 2
	if milieu < prixMin {
		return 0, false
	}
	if round < 3 {
		contreOffre = (milieu + offrePrecedente) / 2
	} else {
		contreOffre = milieu
	}
	return contreOffre, true
}
