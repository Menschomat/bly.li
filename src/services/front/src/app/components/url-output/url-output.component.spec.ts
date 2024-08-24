import { ComponentFixture, TestBed } from '@angular/core/testing';

import { UrlOutputComponent } from './url-output.component';

describe('UrlOutputComponent', () => {
  let component: UrlOutputComponent;
  let fixture: ComponentFixture<UrlOutputComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [UrlOutputComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(UrlOutputComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
