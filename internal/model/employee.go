package model

import "time"

type Employee struct {
	ID                  int64            `json:"ID" db:"id"`
	NamaKaryawan        string           `json:"NAMA_KARYAWAN" db:"name"`
	Jabatan             string           `json:"JABATAN" db:"position"`
	KantorCabang        KantorCabang     `json:"KANTOR_CABANG" db:"branch_office"`
	IsExcluded          int64            `json:"IS_EXCLUDED" db:"is_excluded"`
	GuaranteedDoorprize int64            `json:"GUARANTEED_DOORPRIZE" db:"guaranteed_doorprize"`
	JenisKepegawaian    JenisKepegawaian `json:"JENIS_KEPEGAWAIAN" db:"employment_type"`
	PresentAt           *time.Time       `json:"PRESENT_AT" db:"present_at"`
}

type JenisKepegawaian string

const (
	Organik JenisKepegawaian = "Organik"
	Tad     JenisKepegawaian = "TAD"
)

type KantorCabang string

const (
	KantorCabangBogor        KantorCabang = "Kantor Cabang Bogor"
	KantorCabangCikarang     KantorCabang = "Kantor Cabang Cikarang"
	KantorCabangDepok        KantorCabang = "Kantor Cabang Depok"
	KantorCabangJakartaTimur KantorCabang = "Kantor Cabang Jakarta Timur"
	KantorCabangKarawang     KantorCabang = "Kantor Cabang Karawang"
	KantorCabangPurwokerto   KantorCabang = "Kantor Cabang Purwokerto"
	KantorCabangTangerang    KantorCabang = "Kantor Cabang Tangerang"
	KantorKasBekasi          KantorCabang = "Kantor Kas Bekasi"
	KantorPusatManajemen     KantorCabang = "Kantor Pusat Manajemen"
	KantorPusatOperasional   KantorCabang = "Kantor Pusat Operasional"
)
