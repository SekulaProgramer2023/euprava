import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { UserService, LoginResponse } from '../../../services/user.service2';
import { User } from '../../../model/User'

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [FormsModule, RouterModule],
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']  // svaki komponent ima svoj css
})
export class LoginComponent2 {

  email: string = '';
  password: string = '';
  loginError: string = '';


    user: User = new User(
  '',   // id
  '',   // password
  '',   // role
  '',   // name
  '',   // surname
  '',   // email
  false, // isActive
  [],    // alergije
  []     // omiljenaJela
);


  constructor(private userService: UserService, private router: Router) {}

  onSubmit(): void {
  // Provera username-a i password-a
  if (!this.email || !this.password) {
    this.loginError = 'Please enter both username and password.';
    return;
  }

  const userCredentials = {
    email: this.email,
    password: this.password,
  };

  this.userService.login(this.email, this.password).subscribe({
  next: (res: LoginResponse) => {
    console.log('Login uspešan', res);

    // Sačuvaj token u localStorage
    localStorage.setItem('token', res.token);

    // Preusmeri korisnika
    this.router.navigate(['/menza/home']);
  },
  error: (err) => {
    console.error('Greška pri login-u', err);
    this.loginError = 'Neuspešan login. Proverite email i lozinku.';
  }
});
}



  navigateToRegister() {
    this.router.navigate(['/menza/register']);
  }
}
