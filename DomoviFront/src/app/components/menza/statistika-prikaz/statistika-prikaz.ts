import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FinansijskaKarticaService } from '../../../services/finansijskaKartica.service';

@Component({
  selector: 'app-statistika-prikaz',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './statistika-prikaz.html',
  styleUrl: './statistika-prikaz.css'
})
export class StatistikaPrikaz implements OnInit {
  statistika: any = null;
  loading = true;
  error: string | null = null;

  constructor(private karticaService: FinansijskaKarticaService) {}

  ngOnInit(): void {
    this.karticaService.getStatistika().subscribe({
      next: (res) => {
        this.statistika = res;
        this.loading = false;
      },
      error: (err) => {
        this.error = err.message || 'Greška pri učitavanju statistike';
        this.loading = false;
      }
    });
  }

  getObjectKeys(obj: any): string[] {
    return Object.keys(obj);
  }
}
