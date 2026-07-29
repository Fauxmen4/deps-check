Output docs/demo.gif

Set FontSize 16
Set Width 1200
Set Height 700
Set Theme "Catppuccin Mocha"
Set TypingSpeed 60ms
Set Padding 20

# Быстрый help
Type "depscheck --help"
Sleep 500ms
Enter
Sleep 3s

Ctrl+L

# Основной кейс через SSH URL
Type "depscheck analyze git@github.com:gin-gonic/gin"
Sleep 500ms
Enter
Sleep 5s

Ctrl+L

# Только прямые зависимости
Type "depscheck analyze git@github.com:gin-gonic/gin --direct-only"
Sleep 500ms
Enter
Sleep 4s

Ctrl+L

# JSON формат
Type "depscheck analyze git@github.com:gin-gonic/gin --format json"
Sleep 500ms
Enter
Sleep 5s