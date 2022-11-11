# go-embedings

# 1 : Création d'embedings avec GO
Le programme en Go va créer des embedings pour le jeu de données d'entrée, en tenant de raprocher les elements des mêmes playlists.
La mesure de la performance se fait grace à `cosin_sim` et  `euclidian_dst`.
Cosin_sim : Similarité cosinus : correspondance des vecteurs similaires par multiplication -> plus grande == meilleur
Euclidian_dst : Distance euclidienne entre les points similaires -> plus faible le meilleur.

Ici, on a mesurer la moyenne des metriques pour des pistes similaires et différentes, pour en faire le ratio. On cherche à chaque fois un ratio proche de 0.

```
150100 listes
149693 listes après filtrage
min id: 0 max id: 853777
iter 1 / 13 cosin_sim= 0.9674265706347858 euclidian_dst= 0.5303931167966609 learning_rate= 0.03
iter 2 / 13 cosin_sim= 0.9632334076710268 euclidian_dst= 0.47087193936964744 learning_rate= 0.009
iter 3 / 13 cosin_sim= 0.9526486483550617 euclidian_dst= 0.46574625946734305 learning_rate= 0.0026999999999999997
iter 4 / 13 cosin_sim= 0.9396455937138734 euclidian_dst= 0.46590437513652355 learning_rate= 0.0008099999999999998
iter 5 / 13 cosin_sim= 0.9270189414576175 euclidian_dst= 0.43738314608504136 learning_rate= 0.00024299999999999994
iter 6 / 13 cosin_sim= 0.9213847244260968 euclidian_dst= 0.3856011175212248 learning_rate= 7.289999999999998e-05
iter 7 / 13 cosin_sim= 0.9203596224860883 euclidian_dst= 0.3492285118978613 learning_rate= 2.1869999999999996e-05
iter 8 / 13 cosin_sim= 0.9213495267901786 euclidian_dst= 0.337467267828681 learning_rate= 6.560999999999999e-06
iter 9 / 13 cosin_sim= 0.9223743373181267 euclidian_dst= 0.3353135034821927 learning_rate= 1.9682999999999994e-06
iter 10 / 13 cosin_sim= 0.9227453255049334 euclidian_dst= 0.3350379094142672 learning_rate= 5.904899999999998e-07
iter 11 / 13 cosin_sim= 0.9230000765270144 euclidian_dst= 0.33496533058433603 learning_rate= 1.7714699999999994e-07
iter 12 / 13 cosin_sim= 0.9230672783730042 euclidian_dst= 0.33522758720581025 learning_rate= 5.314409999999998e-08
iter 13 / 13 cosin_sim= 0.9230146893302388 euclidian_dst= 0.33509142217794374 learning_rate= 1.5943229999999993e-08
```

On voit ici que le ratio décroit rapidement pour se stabilisé à 0.33, ce qui signifie que la distance moyenne entre deux pistes d'une même playlist est en moyenne le tier de la distance entre deux pistes pas dans la même playlist.

# 2 : Résultats avec Python
On affiche ensuite les resultats dans un notebook python.
Pour trouver les pistes similaires, on calcul les distances euclidienne entre une piste et toutes les autres.
On a par exemple pour la piste "Thunderstruck de AC/DC" :
-  Back In Black - You Shook Me All Night Long
- L.A. Guns - I Love Rock And Roll
- White Lion - When The Children Cry
- Aerosmith - Water Song/Janie's Got A Gun
- Queensrÿche - Jet City Woman - 2003 Digital Remaster
- AC/DC - Stiff Upper Lip
- Michael Jackson - The Way You Make Me Feel - Single Version
- Saigon Kick - Love Is on the Way
- Meat Loaf - I'd Do Anything For Love (But I Won't Do That) - Longer Still But Not As Long As The Album Version
- AC/DC - The Jack


On peut également prolonger des playlists, avec par exemple une playlist electro (on entree trois musiques de Daft Punk):
- Daft Punk - Around The World
- Daft Punk - Harder Better Faster Stronger
- Daft Punk - Television Rules The Nation / Crescendolls
- Daft Punk - High Fidelity
- Girl Talk - Play Your Part (Pt. 1)
- Pendulum - The Tempest
- Daft Punk - One More Time
- Pharrell Williams - Marilyn Monroe
- Daft Punk - Around The World - Radio Edit
- Benny Benassi presents The Biz - Satisfaction - RL Grime Remix
- Daft Punk - Too Long / Steam Machine
- Daft Punk - Rinzler - Remixed by Kaskade
- Daft Punk - Da Funk
- Justice - Genesis
- C2C - Delta
- The Glitch Mob - Skytoucher
- The Glitch Mob - Animus Vox
- Pharrell Williams - Gust of Wind
- Girl Talk - What It's All About
- Dizzee Rascal - Dance Wiv Me
- Santigold - You'll Find A Way (Switch and Sinden Remix)
- The Glitch Mob - Fortune Days
- M.I.A. - XXXO
- Pretty Lights - Keep Em Bouncin
- Daft Punk - Rectifier
- Daft Punk - Aerodynamic (Daft Punk Remix)