import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ShortTableRowComponent } from './short-table-row.component';

describe('ShortTableRowComponent', () => {
  let component: ShortTableRowComponent;
  let fixture: ComponentFixture<ShortTableRowComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ShortTableRowComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ShortTableRowComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
