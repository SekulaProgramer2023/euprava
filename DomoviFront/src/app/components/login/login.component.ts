import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { UserService } from '../../services/user.service';
import { User } from '../../model/user.model'

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [FormsModule, RouterModule],
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']  // svaki komponent ima svoj css
})
export class LoginComponent {

  email: string = '';
  password: string = '';
  loginError: string = '';


    user: User = new User('', '', '', '', '', '','');

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

  this.userService.login(userCredentials).subscribe({
    next: (res) => {
      console.log('Login uspešan', res);
      localStorage.setItem('user_id', user_id.toString());
      this.router.navigate(['/home']);
    },
    error: (err) => {
      console.error('Greška pri login-u', err);
    }
  });
}


  navigateToRegister() {
    this.router.navigate(['/register']);
  }
}
