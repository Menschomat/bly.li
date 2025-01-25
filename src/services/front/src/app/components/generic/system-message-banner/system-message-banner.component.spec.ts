import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SystemMessageBannerComponent } from './system-message-banner.component';

describe('SystemMessageBannerComponent', () => {
  let component: SystemMessageBannerComponent;
  let fixture: ComponentFixture<SystemMessageBannerComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SystemMessageBannerComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SystemMessageBannerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
