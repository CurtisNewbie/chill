import { ComponentFixture, TestBed } from '@angular/core/testing';

import { BuildInfoComponent } from './build-info.component';

describe('BuildInfoComponent', () => {
  let component: BuildInfoComponent;
  let fixture: ComponentFixture<BuildInfoComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [BuildInfoComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(BuildInfoComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
