import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KreiranjeJela } from './kreiranje-jela';

describe('KreiranjeJela', () => {
  let component: KreiranjeJela;
  let fixture: ComponentFixture<KreiranjeJela>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KreiranjeJela]
    })
    .compileComponents();

    fixture = TestBed.createComponent(KreiranjeJela);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
