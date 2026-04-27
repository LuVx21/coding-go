package fmt_x

import (
	"fmt"
	"io"
	"os"

	au "github.com/logrusorgru/aurora"
)

// xxx: 直接输出
// Fxxx: 指定Writer
// Sxxx: 返回字符串
// xxxf: 格式化输出
// xxxln: 直接输出并换行
// info warn error

var (
	output = os.Stdout
)

type (
	color = func(any) au.Value
)

func Success(a ...any) (int, error)                           { return Fsuccess(output, a...) }
func Successf(ft string, a ...any) (int, error)               { return Fsuccessf(output, ft, a...) }
func Successln(a ...any) (int, error)                         { return Fsuccessln(output, a...) }
func Fsuccess(w io.Writer, a ...any) (int, error)             { return f_xxx(au.Green, w, a...) }
func Fsuccessf(w io.Writer, ft string, a ...any) (int, error) { return f_xxx_f(au.Green, w, ft, a...) }
func Fsuccessln(w io.Writer, a ...any) (int, error)           { return f_xxx_ln(au.Green, w, a...) }
func Ssuccess(a ...any) string                                { return s_xxx(au.Green, a...) }
func Ssuccessf(ft string, a ...any) string                    { return s_xxx_f(au.Green, ft, a...) }
func Ssuccessln(a ...any) string                              { return s_xxx_ln(au.Green, a...) }

// ---------------------------------------------------------------------------------------------------------------------

func Debug(a ...any) (int, error)                           { return Fdebug(output, a...) }
func Debugf(ft string, a ...any) (int, error)               { return Fdebugf(output, ft, a...) }
func Debugln(a ...any) (int, error)                         { return Fdebugln(output, a...) }
func Fdebug(w io.Writer, a ...any) (int, error)             { return f_xxx(au.Blue, w, a...) }
func Fdebugf(w io.Writer, ft string, a ...any) (int, error) { return f_xxx_f(au.Blue, w, ft, a...) }
func Fdebugln(w io.Writer, a ...any) (int, error)           { return f_xxx_ln(au.Blue, w, a...) }
func Sdebug(a ...any) string                                { return s_xxx(au.Blue, a...) }
func Sdebugf(ft string, a ...any) string                    { return s_xxx_f(au.Blue, ft, a...) }
func Sdebugln(a ...any) string                              { return s_xxx_ln(au.Blue, a...) }

// ---------------------------------------------------------------------------------------------------------------------

func Info(a ...any) (int, error)                           { return Finfo(output, a...) }
func Infof(ft string, a ...any) (int, error)               { return Finfof(output, ft, a...) }
func Infoln(a ...any) (int, error)                         { return Finfoln(output, a...) }
func Finfo(w io.Writer, a ...any) (int, error)             { return f_xxx(au.Cyan, w, a...) }
func Finfof(w io.Writer, ft string, a ...any) (int, error) { return f_xxx_f(au.Cyan, w, ft, a...) }
func Finfoln(w io.Writer, a ...any) (int, error)           { return f_xxx_ln(au.Cyan, w, a...) }
func Sinfo(a ...any) string                                { return s_xxx(au.Cyan, a...) }
func Sinfof(ft string, a ...any) string                    { return s_xxx_f(au.Cyan, ft, a...) }
func Sinfoln(a ...any) string                              { return s_xxx_ln(au.Cyan, a...) }

// ---------------------------------------------------------------------------------------------------------------------

func Warning(a ...any) (int, error)                           { return Fwarning(output, a...) }
func Warningf(ft string, a ...any) (int, error)               { return Fwarningf(output, ft, a...) }
func Warningln(a ...any) (int, error)                         { return Fwarningln(output, a...) }
func Fwarning(w io.Writer, a ...any) (int, error)             { return f_xxx(au.Yellow, w, a...) }
func Fwarningf(w io.Writer, ft string, a ...any) (int, error) { return f_xxx_f(au.Yellow, w, ft, a...) }
func Fwarningln(w io.Writer, a ...any) (int, error)           { return f_xxx_ln(au.Yellow, w, a...) }
func Swarning(a ...any) string                                { return s_xxx(au.Yellow, a...) }
func Swarningf(ft string, a ...any) string                    { return s_xxx_f(au.Yellow, ft, a...) }
func Swarningln(a ...any) string                              { return s_xxx_ln(au.Yellow, a...) }

// ---------------------------------------------------------------------------------------------------------------------

func Error(a ...any) (int, error)                           { return Ferror(output, a...) }
func Errorf(ft string, a ...any) (int, error)               { return Ferrorf(output, ft, a...) }
func Errorln(a ...any) (int, error)                         { return Ferrorln(output, a...) }
func Ferror(w io.Writer, a ...any) (int, error)             { return f_xxx(au.Red, w, a...) }
func Ferrorf(w io.Writer, ft string, a ...any) (int, error) { return f_xxx_f(au.Red, w, ft, a...) }
func Ferrorln(w io.Writer, a ...any) (int, error)           { return f_xxx_ln(au.Red, w, a...) }
func Serror(a ...any) string                                { return s_xxx(au.Red, a...) }
func Serrorf(ft string, a ...any) string                    { return s_xxx_f(au.Red, ft, a...) }
func Serrorln(a ...any) string                              { return s_xxx_ln(au.Red, a...) }

// ---------------------------------------------------------------------------------------------------------------------
func f_xxx(c color, w io.Writer, a ...any) (int, error) { return fmt.Fprint(w, c(fmt.Sprint(a...))) }
func f_xxx_f(c color, w io.Writer, ft string, a ...any) (int, error) {
	return fmt.Fprint(w, s_xxx_f(c, ft, a...))
}
func f_xxx_ln(c color, w io.Writer, a ...any) (int, error) {
	return fmt.Fprintln(w, c(fmt.Sprint(a...)))
}

func s_xxx(c color, a ...any) string              { return fmt.Sprint(c(fmt.Sprint(a...))) }
func s_xxx_f(c color, ft string, a ...any) string { return au.Sprintf(c(ft), a...) }
func s_xxx_ln(c color, a ...any) string           { return fmt.Sprintln(c(fmt.Sprint(a...))) }
