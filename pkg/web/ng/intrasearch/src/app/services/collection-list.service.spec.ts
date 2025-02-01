import { TestBed } from '@angular/core/testing';

import { CollectionListService } from './collection-list.service';

describe('CollectionListService', () => {
  let service: CollectionListService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(CollectionListService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
