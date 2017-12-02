## William Kennedy: Behavior Of Channels

https://www.youtube.com/watch?v=zDCKZn4-dck

Logiranje ne smije zaustaviti produkciju. Korištenjem buffered kanala i select naredbe dropaju se paketi iz loggera ako logger prestane biti funkcionalan. Sve je postignuto koristeći Go primitive.

> ...power of concurrency primitives being in the language!

[Primjer](./logger/logger.go) pokazuje nekoliko dobrih patterna go koda:

- factory funkcija za inicijalizaciju
- unutar factory funkcije se pokreću gorutine potrebne za funkcioniranje paketa
- Close funkcija gasi sve gorutine i sinkrono čeka da se pogase
- čist interface prema van, private varijable enkapslirane
- imenovanje funkcija i varijabli (New vs NewLogger)
