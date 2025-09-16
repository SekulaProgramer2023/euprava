import { ComponentFixture, TestBed } from '@angular/core/testing';

import { StatistikaPrikaz } from './statistika-prikaz';

describe('StatistikaPrikaz', () => {
  let component: StatistikaPrikaz;
  let fixture: ComponentFixture<StatistikaPrikaz>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [StatistikaPrikaz]
    })
    .compileComponents();

    fixture = TestBed.createComponent(StatistikaPrikaz);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
