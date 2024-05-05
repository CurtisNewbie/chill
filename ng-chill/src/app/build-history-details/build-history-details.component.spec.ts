import { ComponentFixture, TestBed } from '@angular/core/testing';

import { BuildHistoryDetailsComponent } from './build-history-details.component';

describe('BuildHistoryDetailsComponent', () => {
  let component: BuildHistoryDetailsComponent;
  let fixture: ComponentFixture<BuildHistoryDetailsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [BuildHistoryDetailsComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(BuildHistoryDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
