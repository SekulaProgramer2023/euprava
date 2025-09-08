import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KreiranjeJelovnik } from './kreiranje-jelovnik';

describe('KreiranjeJelovnik', () => {
  let component: KreiranjeJelovnik;
  let fixture: ComponentFixture<KreiranjeJelovnik>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KreiranjeJelovnik]
    })
    .compileComponents();

    fixture = TestBed.createComponent(KreiranjeJelovnik);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
