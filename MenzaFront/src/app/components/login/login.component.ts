import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { UserService } from '../../services/user.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [FormsModule],
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']  // svaki komponent ima svoj css
})
export class LoginComponent {

  email: string = '';
  password: string = '';

  constructor(private userService: UserService, private router: Router) {}

  onSubmit(): void {
    this.userService.login(this.email, this.password).subscribe({
      next: (res) => {
        console.log('Login uspešan', res);
        
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
