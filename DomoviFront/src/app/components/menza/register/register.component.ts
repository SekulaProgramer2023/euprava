import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { UserService, RegisterData } from '../../../services/user.service2';

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [FormsModule],
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.css']  // sopstveni css
})
export class RegisterComponent2 {

  name: string = '';
  surname: string = '';
  email: string = '';
  password: string = '';
  role: string = 'user';
  isActive: boolean = true;

  constructor(private userService: UserService, private router: Router) {}

  onSubmit(): void {
    const data: RegisterData = {
      name: this.name,
      surname: this.surname,
      email: this.email,
      password: this.password,
      role: this.role,
      isActive: this.isActive
    };
    console.log(data)
    this.userService.register(data).subscribe({
      next: (res) => {
        console.log('Korisnik registrovan', res);
        this.router.navigate(['/menza/login']);
      },
      error: (err) => {
        console.error('Greška pri registraciji', err);
      }
    });
  }

  navigateToLogin() {
    this.router.navigate(['/menza/login']);
  }
}
