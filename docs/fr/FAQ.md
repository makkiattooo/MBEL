# MBEL : FAQ (Foire Aux Questions)

### Q : Pourquoi ne pas simplement utiliser du JSON ?
**R :** Le JSON est un format de données, pas un langage. MBEL est un langage de localisation programmable qui supporte nativement les commentaires, les blocs logiques et les métadonnées IA.

### Q : MBEL supporte-t-il les règles de pluriel pour ma langue ?
**R :** Oui. MBEL inclut des règles CLDR pour plus de 15 langues majeures, dont le français, l'anglais, l'allemand, le polonais, le russe, l'espagnol, le chinois et le japonais.

### Q : Quelle est la rapidité du SDK Go ?
**R :** Extrêmement rapide. MBEL compile les fichiers dans un AST optimisé. La résolution est une simple recherche dans une map en O(1).
