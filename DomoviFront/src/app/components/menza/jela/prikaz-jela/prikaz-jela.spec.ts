import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PrikazJela } from './prikaz-jela';

describe('PrikazJela', () => {
  let component: PrikazJela;
  let fixture: ComponentFixture<PrikazJela>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [PrikazJela]
    })
    .compileComponents();

    fixture = TestBed.createComponent(PrikazJela);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
