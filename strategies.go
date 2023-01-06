package main

type StrategieAcheteur func(int, int, int, int) int
type StrategieFournisseur func(int, int, int, int) (int, bool)

func StrategieAcheteurSimple(offre int, prixMax int, round int, aggressivite int) (contreOffre int) {
	difference := offre - prixMax
	if difference < 0 {
		return offre
	}
	contreOffre = prixMax - (difference * aggressivite)
	return contreOffre
}

func StrategieVendeurSimple(offre int, offrePrecedente int, prixMin int, round int) (contreOffre int, possible bool) {
	milieu := (offrePrecedente + offre) / 2
	if milieu < prixMin {
		return 0, false
	}
	if round < MaxRound {
		contreOffre = (milieu + offrePrecedente) / 2
	} else {
		contreOffre = milieu
	}
	return contreOffre, true
}

func StrategieVendeurBizarre(offre int, offrePrecedente int, prixMin int, round int) (contreOffre int, possible bool) {
	return offre + 10, true
}
