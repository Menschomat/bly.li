import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ShortTableComponent } from './short-table.component';

describe('ShortTableComponent', () => {
  let component: ShortTableComponent;
  let fixture: ComponentFixture<ShortTableComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ShortTableComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ShortTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
