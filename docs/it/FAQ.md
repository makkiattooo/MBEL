# MBEL: FAQ (Domande Frequenti)

### D: Perché non usare semplicemente il JSON?
**R:** Il JSON è un formato di dati, non un linguaggio. MBEL è un linguaggio di localizzazione programmabile che supporta nativamente commenti, blocchi logici e metadati IA.

### D: MBEL supporta le regole del plurale per la mia lingua?
**R:** Sì. MBEL include le regole CLDR per oltre 15 lingue principali, tra cui italiano, inglese, tedesco, francese, polacco, russo, spagnolo, cinese e giapponese.

### D: Quanto è veloce l'SDK Go?
**R:** Estremamente veloce. MBEL compila i file in un AST ottimizzato. La risoluzione è una semplice ricerca in una map O(1).
