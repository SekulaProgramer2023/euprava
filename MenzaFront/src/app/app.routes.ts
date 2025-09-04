import { Routes } from '@angular/router';
import { LoginComponent } from './components/login/login.component';
import { RegisterComponent } from './components/register/register.component';
import { PrikazJela } from './components/jela/prikaz-jela/prikaz-jela';
import { KreiranjeJela } from './components/jela/kreiranje-jela/kreiranje-jela';
import { PrikazJelovnik } from './jelovnik/prikaz-jelovnik/prikaz-jelovnik';

export const routes: Routes = [
  { path: 'jelovnik/jela', component: PrikazJela },
   { path: 'jelovnik/kreiranje-jela', component: KreiranjeJela },
      { path: 'jelovnik/prikaz-jelovnika', component: PrikazJelovnik},
  { path: '', redirectTo: '/login', pathMatch: 'full' },
  { path: 'login', component: LoginComponent },
  { path: 'register', component: RegisterComponent },
  { path: '**', redirectTo: '/login' } // fallback ruta
];
