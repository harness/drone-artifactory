package com.example;

import static org.junit.Assert.assertEquals;
import org.junit.Test;

public class UtilTest {
    @Test
    public void testAdd() {
        assertEquals(5, Util.add(2, 3));
    }
}
