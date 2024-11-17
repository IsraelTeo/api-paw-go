package email

import "fmt"

func GetWelcomeEmailSubject() string {
	return "Bienvenido a nuestra familia de Clínica Paw"
}

func GetWelcomeEmailMessage(userName string) string {
	return fmt.Sprintf(`Hola %s,

¡Bienvenido/a a Clínica Paw! Estamos emocionados de que te unas a nuestro equipo, donde nos esforzamos por brindar el mejor cuidado y atención a nuestros pacientes.

Tu llegada es una adición valiosa para nuestro equipo, y estamos seguros de que contribuirás significativamente a nuestra misión de excelencia en el cuidado clínico.

Si necesitas asistencia o tienes alguna pregunta, nuestro equipo está aquí para ayudarte. No dudes en ponerte en contacto con nosotros.

Saludos cordiales,
El equipo de Clínica Paw`, userName)
}
