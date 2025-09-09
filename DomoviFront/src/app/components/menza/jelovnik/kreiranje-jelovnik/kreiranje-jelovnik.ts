import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { JeloService } from './../../../../services/jela.service';
import { JelovnikService } from './../../../../services/jelovnik.service';
import { Jelo } from './../../../../model/Jelo';
import { JelovnikKreiranje } from './../../../../model/JelovnikKreiranje';
import { Router } from '@angular/router';

@Component({
  selector: 'app-kreiranje-jelovnik',
  templateUrl: './kreiranje-jelovnik.html',
  styleUrls: ['./kreiranje-jelovnik.css'],
  standalone: true,
  imports: [CommonModule, FormsModule] // omogućava *ngFor i [(ngModel)]
})
export class KreiranjeJelovnik implements OnInit {
  dorucakJela: Jelo[] = [];
  rucakJela: Jelo[] = [];
  veceraJela: Jelo[] = [];

  izabraniDorucak: string[] = [];
  izabraniRucak: string[] = [];
  izabranaVecera: string[] = [];

  datum: string = '';
  opis: string = '';

  constructor(
    private jeloService: JeloService,
    private jelovnikService: JelovnikService,
    private router: Router 
  ) {}

  ngOnInit(): void {
    this.jeloService.getJelaByTip('dorucak').subscribe(res => this.dorucakJela = res);
    this.jeloService.getJelaByTip('rucak').subscribe(res => this.rucakJela = res);
    this.jeloService.getJelaByTip('vecera').subscribe(res => this.veceraJela = res);
  }
kreirajJelovnik(): void {
  // Napravi objekat sa sva tri tipa jela
  const jelovnikPayload = {
    dorucak: this.izabraniDorucak.length ? this.izabraniDorucak : [], // bar prazan niz
    rucak: this.izabraniRucak.length ? this.izabraniRucak : [],
    vecera: this.izabranaVecera.length ? this.izabranaVecera : [],
    datum: this.datum,
    opis: this.opis
  };

  console.log("Payload koji šaljem na backend:", jelovnikPayload); // za proveru u konzoli

   this.jelovnikService.kreirajJelovnik(jelovnikPayload).subscribe({
      next: res => {
        alert('Jelovnik uspešno kreiran!');
        this.router.navigate(['/menza/jelovnik/prikaz-jelovnika']); // navigacija na prikaz
      },
      error: err => console.error('Greška pri kreiranju jelovnika:', err)
    });
}
}