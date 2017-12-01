Sameer Ajmani: Simulating a real-world system in Go

Postupak pravljenja napitka se može podijeliti u 4 stagea (grindBeans, makeEspresso, steamMilk, makeLatte). 

Mutex lock na cijeli napitak
---

Ako stavimo jedan lock na cijeli postupak pravljenja napitka tada je pravljenje napitka usko grlo. U svakom trenutku se može raditi samo na jednom napitku bez obzira na broj gorutina. Ukupan broj napravljenih napitaka ovisi samo o ukupnom vremenu izvođenja pokusa. Dodavanje gorutina nema utjecaj na konačan ishod jer sve gorutine čekaju na isti zaključani resurs.

```
workers: 1 RESULT: 34664
workers: 2 RESULT: 34783
workers: 3 RESULT: 34326
workers: 4 RESULT: 34537
workers: 5 RESULT: 34455
workers: 6 RESULT: 34237
workers: 7 RESULT: 34396
workers: 8 RESULT: 34064
workers: 9 RESULT: 33561
```

Mutex lock na svaki stage zasebno
---
S obzirom da postupak ima 4 stagea moguće je da 4 gorutine istovremeno rade svaka na svom piću ali na drugom stageu. U ovakvom se okruženju gorutine samoinicijativno preslože u pipeline. Broj pića raste linearno sa brojem dodanih gorutina dok se ne dođe do 4 gorutine. Svaka sljedeća gorutina ima vrlo malen utjecaj na konačan broj proizvedenih napitaka.

```
workers: 1 RESULT: 34276
workers: 2 RESULT: 56473
workers: 3 RESULT: 85185
workers: 4 RESULT: 94028  <-- zasićenje pipelinea
workers: 5 RESULT: 95655
workers: 6 RESULT: 95543
workers: 7 RESULT: 96844
workers: 8 RESULT: 96692
workers: 9 RESULT: 96950
```

Unbuffered channel
---
Za svaki stage smo napravili ulazni kanal (na kojem stage primi task) i izlazni kanal (na kojeg stage pošalje rezultat). Efekt je vrlo sličan pipelineu koji se dogodi sa mutexima. Ipak, rezultat je nešto manji jer svaki stage mora čekati da sljedeći stage preuzme task da bi nastavio raditi. U usporedbi sa kafićem to bi značilo da jedan konobar mora držati napitak u ruci dok je drugi konobar ne preuzme da bi mogao raditi na sljedećem napitku.

Buffered channel
---

Buffered channel ispravlja nedostatak iz proslig zadatka tako što omogući stageu da ostavi rezultat u buffered kanalu i nastavi raditi na sljdećem zadatku. U kafiću bi to bio stol na kojeg konobar može ostaviti neki broj napitaka dok ga drugi konobar ne uzme. 

```
cap: 0 workers: 1 RESULT: 76132  // unbuffered channel, lošije od mutexa
cap: 1 workers: 1 RESULT: 87744
cap: 2 workers: 1 RESULT: 92097
cap: 3 workers: 1 RESULT: 91921
cap: 4 workers: 1 RESULT: 93854
cap: 5 workers: 1 RESULT: 93732
cap: 6 workers: 1 RESULT: 93975
cap: 7 workers: 1 RESULT: 94527
cap: 8 workers: 1 RESULT: 85576
cap: 9 workers: 1 RESULT: 94764  // buffered channel, slično mutexima
```
