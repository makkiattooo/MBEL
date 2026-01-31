# MBEL: FAQ (Najczęściej Zadawane Pytania)

### P: Dlaczego nie po prostu JSON?
**O:** JSON to format danych, a nie język. MBEL to programowalny język lokalizacji. JSON nie wspiera komentarzy, bloków logicznych ani metadanych AI natywnie. MBEL sprawia, że pliki są 10x bardziej czytelne.

### P: Czy MBEL wspiera zasady liczby mnogiej dla mojego języka?
**O:** Tak. MBEL posiada wbudowane reguły CLDR dla ponad 15 głównych języków świata, w tym polskiego, angielskiego, niemieckiego, francuskiego, rosyjskiego, hiszpańskiego, chińskiego i japońskiego.

### P: Jak szybkie jest Go SDK?
**O:** Bardzo szybkie. MBEL kompiluje pliki do optymalnego drzewa AST. Funkcja `T` to szybkie wyszukiwanie w mapie O(1). W testach MBEL jest porównywalny lub szybszy niż `i18next` (JS) czy `gettext`.

### P: Czy mogę użyć MBEL z React/Vue?
**O:** Tak! Użyj komendy `mbel compile`, aby wyeksportować pliki do JSON. Możesz go wczytać do dowolnej biblioteki i18n w JS.

### P: Co się stanie, gdy brakuje klucza?
**O:** Funkcja `T` zwraca sam klucz jako fallback. Dzięki temu aplikacja nie "wybucha" i nie pokazuje pustych miejsc.

### P: Moje tłumaczenie ma `{n}`, a przekazałem string.
**O:** To nie problem. MBEL jest elastyczny. Spróbuje zamienić argument na tekst i wstawi go w miejsce `{n}`.
