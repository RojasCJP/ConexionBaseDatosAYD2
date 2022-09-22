package main

type Usuario struct {
	Username     string `json:"username"`
	Nombre       string `json:"nombre"`
	Correo       string `json:"correo"`
	Contrasena   string `json:"contrasena"`
	Tipo         string `json:"tipo"`         // Administrador, médico, secretaria, cliente
	Estado       int    `json:"estado"`       // Confirmado o no
	Especialidad string `json:"especialidad"` // Especialidad del medico
}

type Mascota struct {
	IdMascota string `json:"idMascota"`
	Nombre    string `json:"nombre"`
	Raza      string `json:"raza"`
	Edad      int    `json:"edad"`
	Foto_Url  string `json:"foto_url"` // URL de S3 de la foto del perro
	Username  string `json:"username"` // Usuario del dueño de la mascota
}

type Sesion struct { // Consulta, cita, emergencia, etc
	IdSesion     string   `json:"idSesion"`
	Creacion     string   `json:"creacion"`     // Hora en que se creó la cita y para validar los 4 minutos para confirmar
	Estado       int      `json:"estado"`       // Pagada o sin confirmar
	Hora_Ingreso string   `json:"hora_ingreso"` // Hora que inició la cita o emergencia
	Hora_Salida  string   `json:"hora_salida"`  // Hora que terminó la cita o emergencia
	Medicamentos []string `json:"medicamentos"` // Arreglo de medicamentos recetados
	Tipo         string   `json:"tipo"`         // Cita o emergencia
	Imagenes     []string `json:"imagenes"`     // Arreglo con las URLS de S3
	IdMascota    string   `json:"idMascota"`    // Id de la mascota atendida
	Username     string   `json:"username"`     // Usuario del médico
}

type Logs struct {
	IdLog          string      `json:"idLog"`
	Metodo         string      `json:"metodo"`
	Entrada        interface{} `json:"entrada"` // Puede ser un objecto con N cantidad de parámetros
	Salida         interface{} `json:"salida"`  // Puede ser un objecto con N cantidad de parámetros
	Error          int         `json:"error"`   // 0 o 1
	Fecha_hora     string      `json:"fecha_hora"`
	Unix_Timestamp int         `json:"unix_timestamp"` // Es la sort key (para ordenamiento)
}

type MascotaImage struct {
	IdMascota    string `json:"idMascota"`
	Nombre       string `json:"nombre"`
	Raza         string `json:"raza"`
	Edad         int    `json:"edad"`
	Foto_Url     string `json:"foto_url"` // URL de S3 de la foto del perro
	Username     string `json:"username"` // Usuario del dueño de la mascota
	ProfilePhoto string `json:"profilephoto"`
}

// Para guardar precios y el descuento
type Precios struct {
	Motivo string  `json:"motivo"` // Medicina general, Ginecología, Descuento, etc
	Precio float64 `json:"precio"`
}
