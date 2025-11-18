package fmt_x

import (
	"fmt"
	"io"

	"github.com/logrusorgru/aurora"
	"github.com/mattn/go-colorable"
)

var (
	output = colorable.NewColorableStdout()
)

// Fsuccessf writes green colored text in manner of fmt.Fprintf
func Fsuccessf(w io.Writer, format string, a ...any) (n int, err error) {
	n, err = fmt.Fprint(w, aurora.Sprintf(aurora.Green(format), a...))
	return
}

// Successf prints green colored text in manner of fmt.Printf
func Successf(format string, a ...any) (n int, err error) {
	n, err = Fsuccessf(output, format, a...)
	return
}

// Ssuccessf returns green colored string in manner of fmt.Sprintf
func Ssuccessf(format string, a ...any) string {
	return aurora.Sprintf(aurora.Green(format), a...)
}

// Finfof writes cyan colored text in manner of fmt.Fprintf
func Finfof(w io.Writer, format string, a ...any) (n int, err error) {
	n, err = fmt.Fprint(w, aurora.Sprintf(aurora.Cyan(format), a...))
	return
}

// Infof prints cyan colored text in manner of fmt.Printf
func Infof(format string, a ...any) (n int, err error) {
	n, err = Finfof(output, format, a...)
	return
}

// Sinfof returns cyan colored string in manner of fmt.Sprintf
func Sinfof(format string, a ...any) string {
	return aurora.Sprintf(aurora.Cyan(format), a...)
}

// Fwarningf writes yellow colored text in manner of fmt.Fprintf
func Fwarningf(w io.Writer, format string, a ...any) (n int, err error) {
	n, err = fmt.Fprint(w, aurora.Sprintf(aurora.Yellow(format), a...))
	return
}

// Warningf prints yellow colored text in manner of fmt.Printf
func Warningf(format string, a ...any) (n int, err error) {
	n, err = Fwarningf(output, format, a...)
	return
}

// Swarningf returns yellow colored string in manner of fmt.Sprintf
func Swarningf(format string, a ...any) string {
	return aurora.Sprintf(aurora.Yellow(format), a...)
}

// Ferrorf writes red colored text in manner of fmt.Fprintf
func Ferrorf(w io.Writer, format string, a ...any) (n int, err error) {
	n, err = fmt.Fprint(w, aurora.Sprintf(aurora.Red(format), a...))
	return
}

// Errorf prints red colored text in manner of fmt.Printf
func Errorf(format string, a ...any) (n int, err error) {
	n, err = Ferrorf(output, format, a...)
	return
}

// Serrorf returns red colored string in manner of fmt.Sprintf
func Serrorf(format string, a ...any) string {
	return aurora.Sprintf(aurora.Red(format), a...)
}

// Fsuccess prints green colored text in manner of fmt.Fprint
func Fsuccess(w io.Writer, a ...any) (n int, err error) {
	n, err = fmt.Fprint(w, aurora.Green(fmt.Sprint(a...)))
	return
}

// Success prints green colored text in manner of fmt.Print
func Success(a ...any) (n int, err error) {
	n, err = Fsuccess(output, a...)
	return
}

// Ssuccess returns green colored string in manner of fmt.Sprint
func Ssuccess(a ...any) string {
	return fmt.Sprint(aurora.Green(fmt.Sprint(a...)))
}

// Finfo prints cyan colored text in manner of fmt.Fprint
func Finfo(w io.Writer, a ...any) (n int, err error) {
	n, err = fmt.Fprint(w, aurora.Cyan(fmt.Sprint(a...)))
	return
}

// Info prints cyan colored text in manner of fmt.Print
func Info(a ...any) (n int, err error) {
	n, err = Finfo(output, a...)
	return
}

// Sinfo returns cyan colored string in manner of fmt.Sprint
func Sinfo(a ...any) string {
	return fmt.Sprint(aurora.Cyan(fmt.Sprint(a...)))
}

// Fwarning prints yellow colored text in manner of fmt.Fprint
func Fwarning(w io.Writer, a ...any) (n int, err error) {
	n, err = fmt.Fprint(w, aurora.Yellow(fmt.Sprint(a...)))
	return
}

// Warning prints yellow colored text in manner of fmt.Print
func Warning(a ...any) (n int, err error) {
	n, err = Fwarning(output, a...)
	return
}

// Swarning returns yellow colored string in manner of fmt.Sprint
func Swarning(a ...any) string {
	return fmt.Sprint(aurora.Yellow(fmt.Sprint(a...)))
}

// Ferror prints red colored text in manner of fmt.Fprint
func Ferror(w io.Writer, a ...any) (n int, err error) {
	n, err = fmt.Fprint(w, aurora.Red(fmt.Sprint(a...)))
	return
}

// Error prints red colored text in manner of fmt.Print
func Error(a ...any) (n int, err error) {
	n, err = Ferror(output, a...)
	return
}

// Serror returns red colored string in manner of fmt.Sprint
func Serror(a ...any) string {
	return fmt.Sprint(aurora.Red(fmt.Sprint(a...)))
}

// Fsuccessln prints green colored text in manner of fmt.Fprintln
func Fsuccessln(w io.Writer, a ...any) (n int, err error) {
	n, err = fmt.Fprintln(w, aurora.Green(fmt.Sprint(a...)))
	return
}

// Successln prints green colored text in manner of fmt.Println
func Successln(a ...any) (n int, err error) {
	n, err = Fsuccessln(output, a...)
	return
}

// Ssuccessln returns green colored string in manner of fmt.Sprintln
func Ssuccessln(a ...any) string {
	return fmt.Sprintln(aurora.Green(fmt.Sprint(a...)))
}

// Finfoln prints cyan colored text in manner of fmt.Fprintln
func Finfoln(w io.Writer, a ...any) (n int, err error) {
	n, err = fmt.Fprintln(w, aurora.Cyan(fmt.Sprint(a...)))
	return
}

// Infoln prints cyan colored text in manner of fmt.Println
func Infoln(a ...any) (n int, err error) {
	n, err = Finfoln(output, a...)
	return
}

// Sinfoln returns cyan colored string in manner of fmt.Sprintln
func Sinfoln(a ...any) string {
	return fmt.Sprintln(aurora.Cyan(fmt.Sprint(a...)))
}

// Fwarningln prints yellow colored text in manner of fmt.Fprintln
func Fwarningln(w io.Writer, a ...any) (n int, err error) {
	n, err = fmt.Fprintln(w, aurora.Yellow(fmt.Sprint(a...)))
	return
}

// Warningln prints yellow colored text in manner of fmt.Println
func Warningln(a ...any) (n int, err error) {
	n, err = Fwarningln(output, a...)
	return
}

// Swarningln returns yellow colored string in manner of fmt.Sprintln
func Swarningln(a ...any) string {
	return fmt.Sprintln(aurora.Yellow(fmt.Sprint(a...)))
}

// Ferrorln prints red colored text in manner of fmt.Fprintln
func Ferrorln(w io.Writer, a ...any) (n int, err error) {
	n, err = fmt.Fprintln(w, aurora.Red(fmt.Sprint(a...)))
	return
}

// Errorln prints red colored text in manner of fmt.Println
func Errorln(a ...any) (n int, err error) {
	n, err = Ferrorln(output, a...)
	return
}

// Serrorln returns red colored string in manner of fmt.Sprintln
func Serrorln(a ...any) string {
	return fmt.Sprintln(aurora.Red(fmt.Sprint(a...)))
}
