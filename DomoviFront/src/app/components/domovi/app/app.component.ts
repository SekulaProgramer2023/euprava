import { Component, OnInit } from '@angular/core';
import { Router, NavigationEnd } from '@angular/router';
import { FaviconService } from '../../../services/favicon,service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {

  constructor(private router: Router, private faviconService: FaviconService) {}

  ngOnInit() {
    this.router.events.subscribe(event => {
      if (event instanceof NavigationEnd) {
        // Promeni favicon u zavisnosti od rute
 if (event.url.includes('home')) {
          this.faviconService.setFavicon('assets/dormitory.png');
        } else {
          this.faviconService.setFavicon('assets/favicon.ico');
        }
      }
    });
  }
}
