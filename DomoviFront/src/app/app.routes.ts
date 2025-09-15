import { Routes } from '@angular/router';
import { LoginComponent } from './components/domovi/login/login.component';
import { RegisterComponent } from './components/domovi/register/register.component';
import { HomeComponent } from './components/domovi/home/home.component'
import { ProfileComponent } from './components/domovi/profile/profile.component';
import { DogadjajComponent } from './components/domovi/dogadjaj/dogadjaj.component';

import { LoginComponent2 } from './components/menza/login/login.component';
import { RegisterComponent2 } from './components/menza/register/register.component';
import { PrikazJela } from './components/menza/jela/prikaz-jela/prikaz-jela';
import { KreiranjeJela } from './components/menza/jela/kreiranje-jela/kreiranje-jela';
import { PrikazJelovnik } from './components/menza/jelovnik/prikaz-jelovnik/prikaz-jelovnik';
import { HomeComponent2 } from './components/menza/home/home.component';
import { ProfileComponent2 } from './components/menza/profile/profile.component';
import { KreiranjeJelovnik } from './components/menza/jelovnik/kreiranje-jelovnik/kreiranje-jelovnik';
import { NotificationComponent } from './components/domovi/notifications/notifications.component'
import { Notifications } from './components/menza/notifications/notifications'
export const routes: Routes = [
  { path: '', redirectTo: '/domovi/login', pathMatch: 'full' },
  { path: 'domovi/login', component: LoginComponent },
  { path: 'domovi/register', component: RegisterComponent },
  { path: 'domovi/home', component: HomeComponent },
  { path: 'domovi/profile', component: ProfileComponent },
  { path: 'domovi/dogadjaj', component: DogadjajComponent },
  { path: 'domovi/notifications', component:  NotificationComponent},


  { path: 'menza/jelovnik/jela', component: PrikazJela },
  { path: 'menza/jelovnik/kreiranje-jela', component: KreiranjeJela },
  { path: 'menza/jelovnik/kreiranje-jelovnika', component: KreiranjeJelovnik },
  { path: 'menza/jelovnik/prikaz-jelovnika', component: PrikazJelovnik},
  { path: '', redirectTo: '/login', pathMatch: 'full' },
  { path: 'menza/login', component: LoginComponent2 },
  { path: 'menza/home', component: HomeComponent2 },
  { path: 'menza/register', component: RegisterComponent2 },
  { path: 'menza/profile', component: ProfileComponent2 },
  { path: 'menza/notifications', component:  Notifications},
];
