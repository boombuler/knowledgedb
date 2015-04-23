# Hilfe

## Suche

### Allgemein
Die Suche, sucht im Standard sowohl nach Titel als auch nach Text.
Dabei wird jedes Dokument als Treffer gewertet, das mindestens einen der Begriffe beinhaltet.

### Zusammenhängende Begriffe

Um zusammenhängende Begriffe zu suchen können diese mittels `"` zusammegefasst werden. 
Beispiel: `"hoher Baum"` findet `hohe Bäume` aber nicht `hohe Häuser und Bäume`.

### Suche in Feldern

Die Suche kann auf bestimmte Felder eingegrenzt werden.
z.B. `Tags:Hilfe` sucht nach Begriffen mit dem Tag `Hilfe`.
Mögliche Felder sind:

* **Tags** sucht in den Tags. Achtet dabei auf den genauen Wortlaut incl. Groß- und Kleinschreibung.
* **Title** sucht im Titel nach dem Wortstamm.
* **Body** sucht im Text nach dem Wortstamm.

### Quantifizierer
Per Prefix kann angegeben werden, ob ein Suchbegriff enthalten sein muß oder nicht enthalten sein darf.

* `Baum` der Text *sollte* Baum enthalten.
* `+Baum` der Text *muß* Baum enthalten
* `-Baum` der Text *darf nicht* Baum enthalten.

### Boosting
Wenn ein Bestimmter Begriff stärker gewichtet werden soll als ein anderer, so kann mittels `^` eine Gewichtung festgelegt werden.
z.B. `Schaufel Wasser^3 Baum^4` sucht Einträge mit den Worten *Schaufel*, *Wasser* und *Baum* wobei Baum die besten Ergebnisse Liefert, *Wasser* die zweitbesten und *Schaufel* die schlechtesten.


### Komplexes Beispiel
`+Title:Baum -Tags:Garten "hoher Baum"^2 Schaufel` würde nach allen Texten suchen,
die im Titel das Wort *Baum* haben, nicht den Tag *Garten* und  *hoher Baum* oder *Schaufel* enthalten,
wobei Texte mit *hoher Baum* weiter oben in den Suchergebnissen stehen würden als welche die nur das Wort *Schaufel* enthalten.
