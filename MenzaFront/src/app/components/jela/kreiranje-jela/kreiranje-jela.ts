import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router'; // <-- dodaj Router
import { JeloService } from '../../../services/jela.service';
import { Jelo, KategorijaJela, TipObroka } from '../../../model/Jelo';

@Component({
  selector: 'app-kreiranje-jela',
  standalone: true,
  imports: [FormsModule, CommonModule],
  templateUrl: './kreiranje-jela.html',
  styleUrls: ['./kreiranje-jela.css']
})
export class KreiranjeJela {
  jelo: Jelo = this.initJelo();

  kategorije = Object.values(KategorijaJela);
  tipoviObroka = Object.values(TipObroka);
  nutritijentiKeys: string[] = ['proteini', 'masti', 'ugljeniHidrati'];

  constructor(private jeloService: JeloService, private router: Router) {} // <-- injektuj Router

  private initJelo(): Jelo {
    return {
      naziv: '',
      kategorija: KategorijaJela.Meso,
      tipObroka: TipObroka.Dorucak,
      kalorije: 0,
      nutritijenti: {
        proteini: 0,
        masti: 0,
        ugljeniHidrati: 0
      }
    };
  }

  addJelo(): void {
    const parsed: { [key: string]: number } = {};
    for (const key of this.nutritijentiKeys) {
      parsed[key] = Number(this.jelo.nutritijenti[key]) || 0;
    }
    this.jelo.nutritijenti = parsed;

    this.jeloService.createJelo(this.jelo).subscribe({
      next: () => {
        alert('Jelo uspešno kreirano!');
        // Navigacija na stranicu sa listom jela
        this.router.navigate(['/jelovnik/jela']);
      },
      error: (err) => {
        console.error(err);
        alert('Greška prilikom kreiranja jela!');
      }
    });
  }
}
