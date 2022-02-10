import unittest

from mask import Mask


class TestMask(unittest.TestCase):

    def test_usage(self):
        m = Mask()
        m.add(1)
        m.add(4)
        m.add(7)

        self.assertEqual(3, len(m))
        self.assertTrue(4 in m)
        self.assertFalse(9 in m)
        self.assertEqual([1, 4, 7], [val for val in m])

        m.pop(4)
        self.assertTrue(7 in m)
        self.assertFalse(4 in m)
        self.assertEqual(2, len(m))
        self.assertEqual([1, 7], [val for val in m])


if __name__ == '__main__':
    unittest.main()
