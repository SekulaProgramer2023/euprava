export class User {
   id: string;
  password: string;
  role: string;
  name: string;
  surname: string;
  email: string;
  soba: string | null;  // <- važno: može biti string ili null
  isActive: boolean;

  constructor(id: string,soba: string, password: string, role: string, name: string, surname: string, email: string) {
    this.id = id;
    this.password = password;
    this.role = role;
    this.name = name;
    this.surname = surname;
    this.soba = soba;
    this.email = email;
    this.isActive = false;
  }
}

