import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SpaceListButtonComponent } from './space-list-button.component';

describe('SpaceListButtonComponent', () => {
  let component: SpaceListButtonComponent;
  let fixture: ComponentFixture<SpaceListButtonComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SpaceListButtonComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SpaceListButtonComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
