import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PrikazJelovnik } from './prikaz-jelovnik';

describe('PrikazJelovnik', () => {
  let component: PrikazJelovnik;
  let fixture: ComponentFixture<PrikazJelovnik>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [PrikazJelovnik]
    })
    .compileComponents();

    fixture = TestBed.createComponent(PrikazJelovnik);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
