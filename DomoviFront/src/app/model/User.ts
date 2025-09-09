export class User {
  id: string;
  password: string;
  role: string;
  name: string;
  surname: string;
  email: string;
  isActive: boolean;
  alergije: string[];
  omiljenaJela: string[];

  constructor(
    id: string,
    password: string,
    role: string,
    name: string,
    surname: string,
    email: string,
    isActive: boolean = false,
    alergije: string[] = [],
    omiljenaJela: string[] = []
  ) {
    this.id = id;
    this.password = password;
    this.role = role;
    this.name = name;
    this.surname = surname;
    this.email = email;
    this.isActive = isActive;
    this.alergije = alergije;
    this.omiljenaJela = omiljenaJela;
  }
}
