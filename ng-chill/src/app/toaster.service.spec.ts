import { TestBed } from '@angular/core/testing';

import { Toaster } from './toaster.service';

describe('Toaster', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: Toaster = TestBed.get(Toaster);
    expect(service).toBeTruthy();
  });
});
