import { Routes } from '@angular/router';
import { LoginComponent } from './components/login/login.component';
import { RegisterComponent } from './components/register/register.component';
import { PrikazJela } from './components/jela/prikaz-jela/prikaz-jela';
import { KreiranjeJela } from './components/jela/kreiranje-jela/kreiranje-jela';
import { PrikazJelovnik } from './jelovnik/prikaz-jelovnik/prikaz-jelovnik';
import { HomeComponent } from './components/home/home.component';
import { ProfileComponent } from './components/profile/profile.component';
import { KreiranjeJelovnik } from './jelovnik/kreiranje-jelovnik/kreiranje-jelovnik';
export const routes: Routes = [
  { path: 'jelovnik/jela', component: PrikazJela },
   { path: 'jelovnik/kreiranje-jela', component: KreiranjeJela },
     { path: 'jelovnik/kreiranje-jelovnika', component: KreiranjeJelovnik },
      { path: 'jelovnik/prikaz-jelovnika', component: PrikazJelovnik},
  { path: '', redirectTo: '/login', pathMatch: 'full' },
  { path: 'login', component: LoginComponent },
   { path: 'home', component: HomeComponent },
  { path: 'register', component: RegisterComponent },
    { path: 'profile', component: ProfileComponent },
  { path: '**', redirectTo: '/login' } // fallback ruta
];
