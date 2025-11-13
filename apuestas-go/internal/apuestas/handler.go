package apuestas

import (
    "database/sql"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"
)

// Handler agrupa las dependencias (como la BD)
type Handler struct {
    DB *sqlx.DB
}

// CrearApuesta es el manejador para el endpoint POST /apuestas
// @Summary      Crea una nueva apuesta
// @Description  Registra una nueva apuesta de un usuario para un evento.
// @Tags         apuestas
// @Accept       json
// @Produce      json
// @Param        apuesta body NuevaApuestaRequest true "Datos de la nueva apuesta"
// @Success      201 {object} Apuesta
// @Failure      400 {object} map[string]string "Error: Datos inválidos"
// @Failure      500 {object} map[string]string "Error: Interno del servidor"
// @Router       /apuestas [post]
func (h *Handler) CrearApuesta(c *gin.Context) {
    var req NuevaApuestaRequest
    
    // 1. Validar el JSON de entrada
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // --- LÓGICA DE SEMANA 3 (Platzhalter) ---
    // 2. Llamar al servicio de Eventos (gRPC) para obtener la cuota.
    //    cuotaActual := 1.5 // (Valor simulado por ahora)
    // 3. Llamar al servicio de Usuarios (gRPC/REST) para verificar saldo.
    //    if !usuarioTieneSaldoSuficiente { ... }
    // --- FIN LÓGICA SEMANA 3 ---
    
    // Por ahora (Semana 2), simulamos la cuota
    cuotaSimulada := 1.75 

    // 4. Insertar en la Base de Datos
    var nuevaApuesta Apuesta
    query := `
        INSERT INTO apuestas (usuario_id, evento_id, monto_apostado, cuota, estado)
        VALUES ($1, $2, $3, $4, 'pendiente')
        RETURNING id, usuario_id, evento_id, monto_apostado, cuota, estado, creado_en`

    err := h.DB.QueryRowx(
        query, 
        req.UsuarioID, 
        req.EventoID, 
        req.MontoApostado, 
        cuotaSimulada,
    ).StructScan(&nuevaApuesta)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear la apuesta"})
        return
    }

    // --- LÓGICA DE SEMANA 3 (Platzhalter) ---
    // 5. Publicar evento en RabbitMQ ("apuesta_creada")
    // --- FIN LÓGICA SEMANA 3 ---

    // 6. Devolver la apuesta creada
    c.JSON(http.StatusCreated, nuevaApuesta)
}
// ObtenerApuesta busca una apuesta por su ID
// @Summary      Obtiene una apuesta por ID
// @Description  Retorna los detalles de una apuesta específica
// @Tags         apuestas
// @Produce      json
// @Param        id   path      string  true  "ID de la Apuesta (UUID)"
// @Success      200  {object}  Apuesta
// @Failure      404  {object}  map[string]string "Error: Apuesta no encontrada"
// @Failure      500  {object}  map[string]string "Error: Interno del servidor"
func (h *Handler) ObtenerApuesta(c *gin.Context) {
    // 1. Obtener el ID de la URL
    id := c.Param("id")

    // 2. Preparar la consulta
    var apuesta Apuesta
    query := "SELECT * FROM apuestas WHERE id = $1"

    // 3. Ejecutar la consulta
    // Usamos h.DB.Get() que es más directo para buscar una sola fila
    err := h.DB.Get(&apuesta, query, id)

    // 4. Manejar errores
    if err != nil {
        if err == sql.ErrNoRows {
            // Error específico: No se encontró la fila
            c.JSON(http.StatusNotFound, gin.H{"error": "Apuesta no encontrada"})
            return
        }
        // Otro error (ej. de conexión a la BD)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar la apuesta"})
        return
    }

    // 5. Devolver la apuesta
    c.JSON(http.StatusOK, apuesta)
}

// ObtenerApuestasPorUsuario busca todas las apuestas de un usuario
// @Summary      Obtiene apuestas por ID de usuario
// @Description  Retorna una lista de todas las apuestas de un usuario
// @Tags         apuestas
// @Produce      json
// @Param        id   path      string  true  "ID del Usuario (UUID)"
// @Success      200  {array}   Apuesta
// @Failure      500  {object}  map[string]string "Error: Interno del servidor"
// @Router       /usuarios/{id}/apuestas [get]
func (h *Handler) ObtenerApuestasPorUsuario(c *gin.Context) {
    // 1. Obtener el ID de usuario de la URL
    usuarioID := c.Param("id")

    // 2. Preparar la consulta
    var apuestas []Apuesta // <-- CAMBIO 1: Debe ser un slice '[]Apuesta'
    query := "SELECT * FROM apuestas WHERE usuario_id = $1 ORDER BY creado_en DESC"

    // 3. Ejecutar la consulta
    // Usamos h.DB.Select() para obtener múltiples filas
    err := h.DB.Select(&apuestas, query, usuarioID) // <-- CAMBIO 2: Usar 'Select'

    // 4. Manejar errores
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar las apuestas"})
        return
    }

    // 5. Devolver las apuestas (devuelve un array vacío si no se encuentran)
    c.JSON(http.StatusOK, apuestas) // <-- CAMBIO 3: Devuelve el slice
}