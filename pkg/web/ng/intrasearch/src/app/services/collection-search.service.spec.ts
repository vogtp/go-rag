import { TestBed } from '@angular/core/testing';

import { CollectionSearchService } from './collection-search.service';

describe('CollectionSearchService', () => {
  let service: CollectionSearchService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(CollectionSearchService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
