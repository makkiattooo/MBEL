# MBEL: FAQ (Häufig gestellte Fragen)

### F: Warum nicht einfach JSON verwenden?
**A:** JSON ist ein Datenformat, keine Sprache. MBEL ist eine programmierbare Lokalisierungssprache. JSON unterstützt nativ keine Kommentare, Logikblöcke oder KI-Metadaten. MBEL macht Ihre Lokalisierungsdateien 10x lesbarer.

### F: Unterstützt MBEL Pluralregeln für meine Sprache?
**A:** Ja. MBEL verfügt über integrierte CLDR-Regeln für über 15 wichtige Weltsprachen, darunter Deutsch, Englisch, Polnisch, Französisch, Russisch, Spanisch, Chinesisch und Japanisch.

### F: Wie schnell ist das Go SDK?
**A:** Extrem schnell. MBEL kompiliert Dateien in einen optimierten AST. Die Auflösung ist ein einfacher O(1) Map-Lookup.

### F: Was passiert, wenn ein Schlüssel fehlt?
**A:** Die Funktion `T` gibt den Schlüssel selbst als Fallback zurück. So wird verhindert, dass Ihre App abstürzt oder leere Stellen anzeigt.
